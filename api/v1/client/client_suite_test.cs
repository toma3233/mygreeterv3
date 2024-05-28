using NUnit.Framework;

namespace MyGreeterV3.Client.Tests
{
    [SetUpFixture]
    public class ClientSuiteSetup
    {
        [OneTimeSetUp]
        public void GlobalSetup()
        {
            // Global initialization logic here
        }

        [OneTimeTearDown]
        public void GlobalTeardown()
        {
            // Global cleanup logic here
        }
    }
}
