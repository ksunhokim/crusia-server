package kim.sunho.crusiaserver
{
	import com.hurlant.crypto.Crypto;
	import com.hurlant.crypto.symmetric.ICipher; 
	import com.hurlant.crypto.symmetric.IPad; 
	import com.hurlant.crypto.symmetric.IVMode; 
	import com.hurlant.crypto.symmetric.PKCS5;
	import com.hurlant.crypto.hash.SHA256;
	import com.hurlant.util.Base64;
	import com.hurlant.util.Hex;
	
	import flash.utils.ByteArray;

	public class Crypto
	{
		static public function sha256(str:String):String {
			var src:ByteArray = Hex.toArray(Hex.fromString(str));
			var sha:SHA256 = new SHA256();
			
			return Hex.fromArray(sha.hash(src));
		}
		
		static public function aes128(key:String, iv:String, str:String):String {
			var inputBA:ByteArray = Hex.toArray(Hex.fromString(str));        
			var keyBA:ByteArray = Hex.toArray(Hex.fromString(key));  
			
			var pad:IPad = new PKCS5();
			var aes:ICipher = com.hurlant.crypto.Crypto.getCipher("aes-128-cbc", keyBA, pad);
			var ivmode:IVMode = aes as IVMode;
			ivmode.IV = Hex.toArray(Hex.fromString(iv));  
			
			aes.encrypt(inputBA); 
			return Base64.encodeByteArray(inputBA);
		}
		
		static public function randomString(strlen:Number):String{
			var chars:String = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
			var num_chars:Number = chars.length - 1;
			var randomChar:String = "";
			
			for (var i:Number = 0; i < strlen; i++){
				randomChar += chars.charAt(Math.floor(Math.random() * num_chars));
			}
			return randomChar;
		}
	}
}