<?php 
namespace Ktools;

class JsonFormat {
    public static function checkFormat(string $jsonStr) {
        try{
            $result = json_decode($jsonStr, true, 512, JSON_UNESCAPED_UNICODE);
            if(is_null($result)) {
                return json_last_error_msg(); 
            } 
            return $jsonStr;
        
        } catch(\JsonException $jsonException) {
            var_export($jsonException);
            return $jsonException->getMessage();
        }
        return 5;
    }

    public static function parser() {
    
    }
}
