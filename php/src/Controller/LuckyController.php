<?php
namespace App\Controller;

use Symfony\Component\HttpFoundation\JsonResponse;

class LuckyController
{
    public function number()
    {
        $relay = new \Spiral\Goridge\SocketRelay("127.0.0.1", 6001);
        $rpc = new \Spiral\Goridge\RPC($relay);

        $t = microtime(true);
        $resp = $rpc->call('custom.Hello', 'world');
        return new JsonResponse([
            'number' => random_int(0, 100),
            'time' => microtime(true) - $t,
            'resp' => $resp,
        ]);
    }
}