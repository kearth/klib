<?php 

namespace Ktools\Json;

class Token {

    private string $type;

    private string $value;

    public function __construct(string $type, $value) {
        $this->type = $type;
        $this->value = $value; 
    }

}
