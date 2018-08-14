package kim.sunho.crusiaserver
{
	import com.hurlant.util.Base64;
	
	import flash.net.URLRequestHeader;
	
	import kim.sunho.crusiaserver.CrusiaServerError;
	import kim.sunho.crusiaserver.Crypto;
	import kim.sunho.crusiaserver.RestClient;
	public class CrusiaServer
	{
		private var url:String;
		private var version:int;
		private var key:String;
		private var iv:String;
		private var token:String;
		
		public function CrusiaServer(url:String, version:int, key:String, iv:String)
		{
			this.url = url;
			this.version = version;
			this.key = Base64.decode(key);
			this.iv = iv;
		}
		
		public function serverVersion(resultHandler:Function, errorHandler:Function):void
		{
			RestClient.execute(url + "/version", "GET", null,
				function(res:Object):void 
				{
					if (!res.data) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT);
					}
					else 
					{
						resultHandler(int(res.data));
					}
				},
				function(status:int, err:String):void
				{
					errorHandler(CrusiaServerError.UNKNOWN, err);
				}
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
					if (!res.data) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT, "");
					}
					else 
					{
						token = res.data;
						resultHandler();
					}	
				},
				function(status:int, err:String):void
				{
					if (err.indexOf("json") >= 0) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT, err);
						return
					}
					switch (status) {
					case 400:
						errorHandler(CrusiaServerError.BAD_FORMAT, err);
						break;
					case 403:
						errorHandler(CrusiaServerError.WRONG_PASSWORD, err);
						break;
					case 404:
						errorHandler(CrusiaServerError.NO_SUCH_USER, err);
						break;
					default:
						errorHandler(CrusiaServerError.UNKNOWN, err);
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
					if (err.indexOf("json") >= 0) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT, err);
						return
					}
					switch (status) 
					{
						case 400:
							errorHandler(CrusiaServerError.BAD_FORMAT, err);
							break;
						case 409:
							errorHandler(CrusiaServerError.EXISTING_USER, err);
							break;
						default:
							errorHandler(CrusiaServerError.UNKNOWN, err);
							break
					}
				}
			);
		}
	
		public function getSaveData(resultHandler:Function, errorHandler:Function):void 
		{
			var header:Array = tokenHeader();
			
			RestClient.execute(url + "/save/get", "POST", {"du":"mmy"},
				function(res:Object):void 
				{
					if (!res.data) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT, "");
					}
					else 
					{
						try 
						{
							resultHandler(JSON.parse(res.data));
						}
						catch(e:TypeError)
						{
							errorHandler(CrusiaServerError.BAD_FORMAT, "");
						}
					}
				},
				function(status:int, err:String):void
				{
					if (err.indexOf("json") >= 0) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT, err);
						return;
					}
					switch (status) 
					{
						case 403:
							errorHandler(CrusiaServerError.UNAUTHORIZED, err);
							break;
						default:
							errorHandler(CrusiaServerError.UNKNOWN, err);
							break
					}
				},
				header,
				20
			);
		}
		
		public function setSaveData(obj:Object, resultHandler:Function, errorHandler:Function):void
		{
			var header:Array = tokenHeader();
			var item:URLRequestHeader = new URLRequestHeader("X-Save-Version", String(version));
			header.push(item);
			var str:String = Crypto.aes128(key, iv, JSON.stringify(obj));
		
			RestClient.execute(url + "/save/set", "POST", str,
				function(res:Object):void 
				{
					resultHandler();
				},
				function(status:int, err:String):void
				{
					if (err.indexOf("json") >= 0) 
					{
						errorHandler(CrusiaServerError.BAD_FORMAT, err);
						return
					}
					switch (status) {
						case 400:
							errorHandler(CrusiaServerError.BAD_FORMAT, err);
							break;
						default:
							errorHandler(CrusiaServerError.UNKNOWN, err);
							break
					}
				},
				header,
				20
			);
		}
		
		private function tokenHeader():Array
		{
			var header:Array = new Array();
			var item:URLRequestHeader = new URLRequestHeader("X-Authorization", token);
			header.push(item);
			return header;
		}
	}
}