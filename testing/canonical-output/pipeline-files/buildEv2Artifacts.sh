# Create Artifacts Files
# Assumptions for the script:
# - If isService is true, 
#   - the directory being used to build artifacts is $directoryName/server
#   - there is a testServiceImage.sh file in the $directoryName/server/test being used to test the service.
#   - there is a deployments folder in the $directoryName/server directory that contains the helm chart for the service.
# - There is an folder labelled "Ev2" in the directory being used to build artifacts
# - There is a resources folder in the directory being used to build artifacts that stores the bicep files for the resources and any corresponding parameters files (starting with "template-")

set -e
directoryName=$1
outputDir=$2
isService=$3
rolloutInfra=$4
buildNumber=$5
isLocal=$6

if [ "$isService" = "true" ]; then
  directory=$directoryName/server
else 
  directory=$directoryName
fi

cd $directory
currPath=$(pwd)
echo "Copy Ev2 folder to out directory"
cp -rT Ev2 $outputDir

echo "Test and Build Directory Specific Code"
if [ "$isService" = "true" ]; then
  cd ..
  cd ..
  echo "Test Service"
  ./${directory}/test/testServiceImage.sh
  
  cd $directory

  mkdir -p $outputDir/Ev2Specs/Build

  if [ "$isLocal" = "true" ]; then
    echo "Building Docker Image"
    cd generated
    docker build --build-arg PAT=$READPAT -t $directoryName -f ../Dockerfile ./../
    docker save -o $directoryName-image.tar $directoryName
    cp $directoryName-image.tar $outputDir/Ev2Specs/Build
    cd ..
    
  fi

  # Run customized build script for the service if it does exist.
  if [ -f "buildCustomEv2.sh" ]; then
    echo "Running build script for service"
    ./buildEv2.sh
  fi

  echo "Package Helm"
  # Install helm if not already installed
  if ! command -v helm &> /dev/null; then
      echo "Helm not found. Installing Helm..."
      curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
      chmod 700 get_helm.sh
      ./get_helm.sh
  fi
  # Package service
  cd deployments
  helm package .
  cd ..

  echo "Copy Helm Package"
  mv deployments/$directoryName-0.1.0.tgz $outputDir/Ev2Specs/Build

  echo "Copy Helm Values"
  find . -name "*values*.yaml" -exec cp {} $outputDir/Ev2Specs/Build \;

  echo "Rename Helm Values files"
  for file in $outputDir/Ev2Specs/Build/template-*.yaml; do
    mv "$file" "${file/template-/}"
  done

  cd $outputDir/Ev2Specs

  # Downloading Crane here so that the Ev2 shell has what it needs to push to the ACR, putting it in with script
  if [ ! -f Shell/crane.tar.gz ]; then
    curl -L -o Shell/crane.tar.gz https://github.com/google/go-containerregistry/releases/download/v0.4.0/go-containerregistry_Linux_x86_64.tar.gz
  fi

  # Change to the Shell directory
  cd Shell

  # Extract the tar.gz file
  tar xzvf crane.tar.gz

  # Return to the previous directory
  cd ..

  # Packaging everything together: script and crane
  tar -cvf push-image-to-acr.tar ./Shell/*
else
  echo "Test Shared Resources"
  ./resources/testResourceNames.sh ev2
fi

cd $currPath

mkdir -p $outputDir/Ev2Specs/Templates
mkdir -p $outputDir/Ev2Specs/Parameters

echo "Copy all bicep files in $directory/resources to $outputDir/Ev2Specs/Templates"
cd resources
find . -name "*.bicep" -exec cp {} $outputDir/Ev2Specs/Templates \;

echo "Copy parameters file to $outputDir/Ev2Specs/Parameters"
find . -name "template-*.json" -exec cp {} $outputDir/Ev2Specs/Parameters \;

echo "Rename parameters file"
for file in $outputDir/Ev2Specs/Parameters/template-*.json; do
  mv "$file" "${file/template-/}"
done

echo "Convert Bicep Templates to json"
cd $outputDir/Ev2Specs/Templates
for f in *.bicep; do az bicep build --file "$f"; done
cd ..

echo "Package Script and Set Build Version"
versionContent=$buildNumber
versionFileName="./Version.txt"

# If the file already exists, delete it so you can recreate it and repopulate with the build number
if [ -f "$versionFileName" ]; then
  rm "$versionFileName"
fi

# Create the version file with the build number
echo -n "$versionContent" > "$versionFileName"

echo "Copy configuration.json"
cp Configurations/$rolloutInfra/Configuration.json .

