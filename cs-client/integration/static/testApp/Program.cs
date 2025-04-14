using System.Net;
using Nicechain;

static byte[] StringToByteArray(string hex)
{
    return Enumerable.Range(0, hex.Length)
                     .Where(x => x % 2 == 0)
                     .Select(x => Convert.ToByte(hex.Substring(x, 2), 16))
                     .ToArray();
}

byte[] privateKey = StringToByteArray(args[0]);
String grpcURL = args[1];

var queryClient = new QueryClient(grpcURL);
var txClient = new TxClient(queryClient, privateKey);

await txClient.NicechainV1TxClient.SendMsgCreateTestItem(
    new Nicechain.Nicechain.V1.MsgCreateTestItem
    {
        Creator = txClient.Ec.AccoutAddress.Bech32
    }
);

Nicechain.Nicechain.V1.QueryAllTestItemResponse response = queryClient.NicechainV1QueryClient.ListTestItem(new Nicechain.Nicechain.V1.QueryAllTestItemRequest { });

Console.Out.WriteLine(response.ToString());

int itemCount = response.TestItem.Count;
if (itemCount != 1)
{
    throw new Exception("expected one TestItem in chain, but got: " + itemCount);
}
