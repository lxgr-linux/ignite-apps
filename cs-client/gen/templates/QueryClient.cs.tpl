using Grpc.Net.Client;

namespace {{ .NameSpace }}
{
    public class QueryClient : Cosmcs.Client.QueryClient
    {
        {{- range .Services }}
        public {{ .Path }}.{{ .Type }}.{{ .Type }}Client {{ .Name }}QueryClient { get; }
        {{- end }}
        public QueryClient(string rpcUrl, GrpcChannelOptions? options = null) : base(rpcUrl, options)
        {
            {{- range .Services }}
            {{ .Name }}{{ .Type }}Client = new {{ .Path }}.{{ .Type }}.{{ .Type }}Client(Channel);
            {{- end }}
        }
    }
}
