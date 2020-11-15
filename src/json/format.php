<?php 

namespace Ktools\Json;

class Format implements \JsonSerializable{
    private array $arr = array();

    public function __construct(array $arr) {
        $this->arr = $arr; 
    } 

    public function jsonSerialize() {
        return $this->arr;
    }
}
