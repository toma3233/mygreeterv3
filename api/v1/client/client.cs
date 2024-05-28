using System;
using Grpc.Net.Client;
using MyGreeterV3.Api.V1;

namespace MyGreeterV3.Client
{
    public class Client
    {
        private MyGreeter.MyGreeterClient _client;

        public Client(string address)
        {
            var channel = GrpcChannel.ForAddress(address);
            _client = new MyGreeter.MyGreeterClient(channel);
        }

        public string SayHello(string name)
        {
            var request = new HelloRequest { Name = name };
            var reply = _client.SayHello(request);
            return reply.Message;
        }
    }
}
