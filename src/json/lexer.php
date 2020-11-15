<?php 

namespace Ktools\Json;

class Lexer {

    private int $index;

    private array $jsonArr;

    private string $json;

    private int $status;

    public function __construct(&array $jsonArr, string $json) {
        $this->json    = $json; 
        $this->jsonArr = $jsonArr;
        $this->status  = 0;
    }    

    public function getToken(int $index) : int {
        

    }

    private function parseObject(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function parseArray(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function parseNull(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function parseFalse(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function parseTrue(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function parseString(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    private function parseNumber(&$arr, $index) : int {
        $index += 1;
        return $index;
    }

    public function nextToken() : object {
    
        switch ($chr) {
        case "{":       
            $this->index = $this->parseObject($this->jsonArr, $this->index);
            break;    
        case "[":        
            $this->index = $this->parseArray($this->jsonArr, $this->index);
            break;    
        case 'n':
            $this->index = $this->parseNull($this->jsonArr, $this->index);
            break;    
        case 't':
            $this->index = $this->parseTrue($this->jsonArr, $this->index);
            break;    
        case 'f':
            $this->index = $this->parseFalse($this->jsonArr, $this->index);
            break;    
        case '"':
            $this->index = $this->parseString($this->jsonArr, $this->index);
            break;    
        default:
            $this->index = $this->parseNumber($this->jsonArr, $this->index);
        }
    }

    public function hasMore() : bool {
    
    }
}
