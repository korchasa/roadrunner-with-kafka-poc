<?php
namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;

class DefaultController
{
    public function healthz()
    {
        return new JsonResponse(["ok" => true]);
    }

    public function number()
    {
        $relay = new \Spiral\Goridge\SocketRelay("127.0.0.1", 6001);
        $rpc = new \Spiral\Goridge\RPC($relay);

        // payload emulation
        $t = microtime(true);
        usleep(50000);
        $payload_time = microtime(true) - $t;

        $msg = json_encode([
            "service" => "router",
            "time_str"=>"2019-04-24T16:11:55.247Z",
            "timestamp"=>1556122315,
            "timestamp_micro"=>247000,
            "host"=>"router-deployment-955b9fd77-jjrlg",
            "type"=>"error",
            "message"=>"Error: InvalidMessage in file `/app/node_modules/kafka-node/lib/protocol/protocol.js` line 764",
            "trace_id"=>"e01a030da6a4fa59:7e17ce4a94b71882:e01a030da6a4fa59:1",
            "uri"=>"/html5/",
            "user_source"=>"unknown",
            "context"=>[],
            "action_params"=>[],
            "user_client"=>"unknown",
            "user_browser"=>"unknown",
            "user_browser_version"=>"unknown",
            "user_operating_system"=>"unknown",
            "user_operating_system_version"=>"unknown",
            "user_device_category"=>"unknown",
            "user_device_brand"=>"unknown",
            "user_device_model"=>"unknown",
            "error_type"=>"Error",
            "caller"=>"undefined,764",
            "error_trace"=>"Error: InvalidMessage
        at Object.<anonymous> (/app/node_modules/kafka-node/lib/protocol/protocol.js:764:17)
        at Object.self.tap (/app/node_modules/binary/index.js:248:12)
        at Object.decodePartitions (/app/node_modules/kafka-node/lib/protocol/protocol.js:762:73)
        at Object.self.loop (/app/node_modules/binary/index.js:267:16)
        at Object.<anonymous> (/app/node_modules/kafka-node/lib/protocol/protocol.js:163:8)
        at Object.self.loop (/app/node_modules/binary/index.js:267:16)
        at decodeProduceV1Response (/app/node_modules/kafka-node/lib/protocol/protocol.js:756:6)
        at KafkaClient.Client.invokeResponseCallback (/app/node_modules/kafka-node/lib/client.js:823:18)
        at KafkaClient.Client.handleReceivedData (/app/node_modules/kafka-node/lib/client.js:803:10)
        at KafkaClient.handleReceivedData (/app/node_modules/kafka-node/lib/kafkaClient.js:1226:48)",
            "error_message"=>"InvalidMessage"
        ]);

        $t = microtime(true);
        $resp = $rpc->call('kafka.Produce', $msg);
        $kafka_time = microtime(true) - $t;
        return new JsonResponse([
            'number' => random_int(0, 100),
            'payload_time' => $payload_time,
            'kafka_time' => $kafka_time,
            'resp' => $resp,
        ]);
    }
}