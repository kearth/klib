<?php 
declare(strict_types=1);

use PHPUnit\Framework\TestCase;

final class JSONFormatTest extends TestCase {
    public function testCheckFormat() : void {
        $json = '{"id":"9467","sid":"7127","lid":"1497","ltype":"1","status":"0","update_time":"1604460903","update_user":"tangxindi","origin_id":"6090","origin_note":"android","is_new":"1","sdkType":"0"}';
        $res = Ktools\JsonFormat::checkFormat($json);
        echo $res;
    }

}
