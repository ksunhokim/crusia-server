package kim.sunho.crusiaserver
{
	import com.hurlant.crypto.hash.SHA256;
	import com.hurlant.util.Hex;
	
	import flash.utils.ByteArray;

	public class Crypto
	{
		static public function sha256(str:String):String {
			var src:ByteArray = new ByteArray;
			src.writeUTFBytes(str);
			var sha:SHA256 = new SHA256();
			return Hex.fromArray(sha.hash(src));
		}
	}
}