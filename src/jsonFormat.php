<?php 
namespace Ktools;

class JsonFormat {

    public static function printJson(array $jsonArr) : string {
        return json_encode($jsonArr, JSON_PRETTY_PRINT);
    }

    public static function checkFormat(string $jsonStr) {
        try{
            $paser = new Json\Parser($jsonStr);
            $paser->parse();
        

        } catch(\JsonException $jsonException) {
            var_export($jsonException);
            return $jsonException->getMessage();
        }
        return 5;
    }

    public static function parser() {
    
    }
}
