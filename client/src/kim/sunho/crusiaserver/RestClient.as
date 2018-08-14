package kim.sunho.crusiaserver
{
	import flash.events.Event;
	import flash.events.IOErrorEvent;
	import flash.net.URLLoader;
	import flash.net.URLLoaderDataFormat;
	import flash.net.URLRequest;

	public class RestClient
	{
		static private var list:Vector.<RestClient> = new Vector.<RestClient>;
		static public function execute(url: String, method:String, params:*, resultHandler:Function, errorHandler:Function, header:Array = null, retries:int = 0):void
		{
			var client:RestClient = new RestClient;
			client.url = url;
			client.method = method;
			client.params = params;
			client.resultHandler = resultHandler;
			client.errorHandler = errorHandler;
			client.header = header;
			client.retries = retries;
			
			list.push(client);
			client.setup();
			client.run();
		}

		private var url:String;
		private var method:String;
		private var params:*;
		private var resultHandler:Function;
		private var errorHandler:Function;
		private var json:Boolean;
		private var header:Array;
		
		private var loader:URLLoader;
		private var retries:int;
		
		public function setup():void
		{
			loader = new URLLoader;
			loader.dataFormat = URLLoaderDataFormat.TEXT;
			loader.addEventListener(Event.COMPLETE, onComplete);
			loader.addEventListener(IOErrorEvent.IO_ERROR, onError);
		}
		
		public function run():void
		{
			var req:URLRequest = new URLRequest(url);
			req.method = method;
			
			if (params is String) 
			{
				req.contentType = "plain/text";
				req.data = params;
			} 
			else if(params)
			{
				req.contentType = "application/json";
				req.data = JSON.stringify(params);
			}
			
			if(header) {
				req.requestHeaders = req.requestHeaders.concat(header);
			}

			loader.load(req);
		}
		
		private function onComplete(e:Event):void
		{
			if (e.target != loader) return;
			
			try 
			{
				var o:Object = JSON.parse(loader.data);
				if (o.status != 200)
				{
					errorHandler(o.status, o.msg);
				}
				else
				{
					resultHandler(o);
				}
			}
			catch(e:TypeError)
			{
				errorHandler(400, "json parse error");
			}
			
			destroy();
		}
		
		private function onError(e:IOErrorEvent):void {
			if (e.target != loader) return;
			
			if (retries != 0) 
			{
				retries --;
				run();
				return;
			}
			
			errorHandler(500, e.text);
			destroy();
		}
		
		
		private function destroy():void 
		{
			var anIndex:int = list.indexOf(this);
			if (anIndex > -1) list.splice(anIndex, 1);
			
			resultHandler = null;
			errorHandler = null;
			params = null;
			
			if (!loader) return;
			
			loader.removeEventListener(Event.COMPLETE, onComplete);
			loader.removeEventListener(Event.COMPLETE, onError);
			loader = null;
		}
	}
}