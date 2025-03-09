using System.Net.Http;
using Cosmcs.Client;

namespace {{ .NameSpace }}
{
    public class Client
    {
        public EasyClient Ec { get; }
        public QueryClient QueryClient { get; }
        {{- range .Services }}
        public {{ .Path }}.{{ .Type }}Client {{ .Name }}TxClient { get; }
        {{- end }}

        public Client(QueryClient queryClient, string chainId, byte[] bytes)
        {
            Ec = new EasyClient( queryClient, chainId, bytes, "cosmos");
            QueryClient = queryClient;
            {{- range .Services }}
            {{ .Name }}TxClient = new {{ .Path }}.{{ .Type }}Client(Ec);
            {{- end }}
        }
    }
}
