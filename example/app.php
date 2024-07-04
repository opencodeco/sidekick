<?php

declare(strict_types=1);

use Hyperf\Nano\Factory\AppFactory;

require_once __DIR__ . '/vendor/autoload.php';

$app = AppFactory::createBase(host: '0.0.0.0', port: 9501);

$app->get('/', function () {
    return 'Hello, World!';
});

$app->run();
