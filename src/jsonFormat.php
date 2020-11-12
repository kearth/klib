<?php 
namespace Ktools

class JsonFormat {
    public static function checkFormat(string $jsonStr) {
        $result = json_decode($jsonStr);
        return $result;
    }
}
