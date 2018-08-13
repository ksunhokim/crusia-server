package kim.sunho.crusiaserver
{
	import com.hurlant.crypto.hash.SHA256;
	import com.hurlant.util.Base64;
	
	import flash.net.URLRequestHeader;
	import flash.utils.ByteArray;
	
	import kim.sunho.crusiaserver.CrusiaServerError;
	import kim.sunho.crusiaserver.Crypto;
	import kim.sunho.crusiaserver.RestClient;
	public class CrusiaServer
	{
		private var url:String;
		private var version:int;
		private var saveKey:String;
		private var token:String;
		
		public function CrusiaServer(url:String, version:int, saveKey:String)
		{
			this.url = url;
			this.version = version;
			this.saveKey = Base64.decode(saveKey);
		}
		
		public function serverVersion(resultHandler:Function, errorHandler:Function):void
		{
			RestClient.execute(url + "/version", "GET", null,
				function(res:Object):void 
				{
					if (res.version) 
					{
						resultHandler(int(res.version));
					}
					else 
					{
						errorHandler(CrusiaServerError.NETWORK);
					}
				},
				function(status:int, err:String):void
				{
					errorHandler(CrusiaServerError.UNKNOWN);
				},
				true
			);
		}
		
		public function login(username:String, password:String, resultHandler:Function, errorHandler:Function):void
		{
			var pwhash:String = Crypto.sha256(password);
			var params:Object = 
			{
				"username": username,
				"passhash": pwhash
			};
			
			RestClient.execute(url + "/login", "POST", params,
				function(res:Object):void 
				{
					if (res.token) 
					{
						token = res.token;
						trace(token);
						resultHandler();
					} 
					else 
					{
						errorHandler(CrusiaServerError.NETWORK);
					}
				},
				function(status:int, err:String):void
				{
					switch (status) {
					case 400:
						errorHandler(CrusiaServerError.BAD_FORMAT);
						break;
					case 403:
						errorHandler(CrusiaServerError.WRONG_PASSWORD);
						break;
					case 404:
						errorHandler(CrusiaServerError.NO_SUCH_USER);
						break;
					default:
						errorHandler(CrusiaServerError.UNKNOWN);
						break
					}
				}
			);
		}
		
		public function register(username:String, password:String, email:String, resultHandler:Function, errorHandler:Function):void
		{
			var pwhash:String = Crypto.sha256(password);
			var params:Object = 
			{
				"username": username,
				"passhash": pwhash,
				"email": email
			};
			
			RestClient.execute(url + "/register", "POST", params,
				function(res:Object):void 
				{
					resultHandler();
				},
				function(status:int, err:String):void
				{
					switch (status) {
						case 400:
							errorHandler(CrusiaServerError.BAD_FORMAT);
							break;
						case 409:
							errorHandler(CrusiaServerError.EXISTING_USER);
							break;
						default:
							errorHandler(CrusiaServerError.UNKNOWN);
							break
					}
				},
				false
			);
		}
	
		public function getSaveData(resultHandler:Function, errorHandler:Function):void 
		{
			var header:Array = tokenHeader();
			
			RestClient.execute(url + "/save/get", "POST", " ",
				function(res:Object):void 
				{
					resultHandler(res);
				},
				function(status:int, err:String):void
				{
					switch (status) {
						case 200:
							if (err.indexOf("json") >= 0) {
								errorHandler(CrusiaServerError.BAD_FORMAT);
							}
							break;
						case 400:
							errorHandler(CrusiaServerError.BAD_FORMAT);
							break;
						case 409:
							errorHandler(CrusiaServerError.EXISTING_USER);
							break;
						default:
							errorHandler(CrusiaServerError.UNKNOWN);
							break
					}
				},
				true,
				header
			);
		}
		
		public function setSaveData(obj:Object, resultHandler:Function, errorHandler:Function):void
		{
			var header:Array = tokenHeader();
			var item:URLRequestHeader = new URLRequestHeader("X-Save-Version", String(version));
			header.push(item);
			trace(JSON.stringify(header));
			
			var str:String = Crypto.aes128(saveKey, JSON.stringify(obj));
			trace(str);
			RestClient.execute(url + "/save/set", "POST", str,
				function(res:String):void 
				{
					resultHandler();
				},
				function(status:int, err:String):void
				{
					switch (status) {
						case 400:
							errorHandler(CrusiaServerError.BAD_FORMAT);
							break;
						case 409:
							errorHandler(CrusiaServerError.EXISTING_USER);
							break;
						default:
							errorHandler(CrusiaServerError.UNKNOWN);
							break
					}
				},
				false,
				header
			);
		}
		
		private function tokenHeader():Array
		{
			var header:Array = new Array();
			var item:URLRequestHeader = new URLRequestHeader("Authorization", token);
			header.push(item);
			return header;
		}
	}
}