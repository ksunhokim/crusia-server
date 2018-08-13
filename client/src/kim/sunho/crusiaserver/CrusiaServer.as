package kim.sunho.crusiaserver
{
	import com.hurlant.crypto.hash.SHA256;
	import com.hurlant.util.Hex;
	
	import flash.utils.ByteArray;
	
	import kim.sunho.crusiaserver.CrusiaServerError;
	import kim.sunho.crusiaserver.RestClient;
	import kim.sunho.crusiaserver.Crypto;
	public class CrusiaServer
	{
		private var url:String;
		private var version:int;
		private var saveKey:String;
		private var token:String;
		
		public function CrusiaServer(url:String, version:int, saveKey:String) {
			this.url = url;
			this.version = version;
			this.saveKey = saveKey;
		}
		
		public function serverVersion(resultHandler:Function, errorHandler:Function):void {
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
		
		public function register(username:String, password:String, email:String, resultHandler:Function, errorHandler:Function):void {
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
	}
}