using NUnit.Framework;
using MyGreeterV3.Client;

namespace MyGreeterV3.Client.Tests
{
    [TestFixture]
    public class ClientTests
    {
        private Client client;

        [SetUp]
        public void Setup()
        {
            // Initialize your client here
            client = new Client("localhost:50051");
        }

        [Test]
        public void TestClientConnection()
        {
            // Test client connection or functionality here
            Assert.IsNotNull(client);
        }
    }
}
