#!/bin/bash -x

# Execute this file under the service directory.

# `go work use`` must be after the module's go.mod exists.
# Once the go.work file exist, the module where you want to run `go build ./...`
# must be in go.work's use list.

# For the api.
cd api
go mod init <<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>
go mod edit -require github.com/Azure/aks-middleware@v0.0.23
go get google.golang.org/genproto@latest
cd ..
cd api
cd v1 
mkdir mock # Current workaround until we find a cleaner way to add the directory to the template files.
make service
if [ $? -ne 0 ]
then
    echo "Make service failed."
    exit 1
fi
cd ..
go build ./...
if [ $? -ne 0 ]
then
    echo "Building the api module failed."
    exit 1
fi
go test ./...
if [ $? -ne 0 ]
then
    echo "Testing the api module failed."
    exit 1
fi
gofmt -s -w .
cd ..
<<- if and (not (contains .envInformation.goModuleNamePrefix "service_hub")) (not (contains .envInformation.goModuleNamePrefix "go.goms.io"))>>
# This process automates tagging a new version of api
# if there were changes from the last tag.
# Increase tag version or set to 0.0.1 if never tagged before.
currentTag=`git describe --abbrev=0 --tags`
currentTag=${currentTag:-'0.0.0'}
version=${currentTag//v/}  # Remove the leading 'v'
version=${version##*/}  # Extract the part after the last '/'
major="${version%%.*}"
minorPatch="${version#*.}"
minor="${minorPatch%%.*}"
patch="${minorPatch#*.}"
patch=$((patch+1))
newTag=$major.$minor.$patch
git diff --quiet $currentTag -- api
if [ $? -ne 0 ]
then
    git add ./api
    git commit -m "api module v$newTag"
    git tag <<.serviceInput.directoryName>>/api/v$newTag
    git push --set-upstream origin master
    git push origin <<.serviceInput.directoryName>>/api/v$newTag
fi
<<end>>

cat <<`<<`>> EOM

If your goModuleNamePrefix has . (dot) inside, you have to create the module
in the repo. Otherwise Go will complain that the module cannot be found
even if you use go.work.

Use the following commands:

git add ./api
git commit -m "api module v0.0.1"
git tag <<.serviceInput.directoryName>>/api/v0.0.1
git push
git push origin <<.serviceInput.directoryName>>/api/v0.0.1

Then you come back here to run this script again.
After git push, the "module cannot be found" message may still persist.
Wait a couple of minutes and the git repo will be able to return the module.

EOM

echo -----------------------

# For the server
cd server
go mod init <<serverModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>
<<- if (contains .envInformation.goModuleNamePrefix "go.goms.io")>>
go mod edit -require <<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>@v0.0.7<<end>>
<<- if (contains .envInformation.goModuleNamePrefix "service_hub")>>
go mod edit -require <<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>@v0.0.7<<end>>
<<- if (contains .envInformation.goModuleNamePrefix "temp_userenv")>>
go mod edit -require <<apiModule .envInformation.goModuleNamePrefix .serviceInput.directoryName>>@v0.0.7<<end>>
go mod edit -require github.com/Azure/aks-middleware@v0.0.23
go mod tidy
# The following command must be run AFTER go mod tidy. If ran before, building the server module
# will fail as go mod tidy removes the indirect dependency with google.golang.org/genproto
# and go work sum will pull in an older version that causes an ambiguous import error.
# For more information refer to: https://github.com/googleapis/go-genproto/issues/1015
go get google.golang.org/genproto@latest
cd ..

go work init
go work use ./api
go work use ./server

cd server
go build ./...
if [ $? -ne 0 ]
then
    echo "Building the server module failed."
    echo "If downloading the server module failed, you might have to wait for the api module to be available or the tag to settle then rerun again"
    exit 1
fi
go test ./...
if [ $? -ne 0 ]
then
    echo "Testing the server module failed."
    exit 1
fi
gofmt -s -w .
cd ..

cat <<`<<`>> EOM

After the service can run properly on your local machine, you can use the commands in
the Makefile in this directory to run the service on your standalone.

!!! Rename/delete your go.work file as aksbuilder doesn't work with go.work.

Remember to commit your modules to the repo and use the right version of the module.
Local changes won't be used by aksbuilder.

EOM
