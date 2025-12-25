export namespace main {
	
	export class ChatWithAIResult {
	    message: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatWithAIResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	        this.error = source["error"];
	    }
	}
	export class GetCurrentUserResult {
	    userId: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new GetCurrentUserResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.error = source["error"];
	    }
	}
	export class Mesaage {
	    id: string;
	    userId: string;
	    content: string;
	    role: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new Mesaage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.userId = source["userId"];
	        this.content = source["content"];
	        this.role = source["role"];
	        this.createdAt = source["createdAt"];
	    }
	}
	export class LoadChatHistoryResult {
	    messages: Mesaage[];
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new LoadChatHistoryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.messages = this.convertValues(source["messages"], Mesaage);
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SignOutResult {
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new SignOutResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.error = source["error"];
	    }
	}
	export class SigninResult {
	    token: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new SigninResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.error = source["error"];
	    }
	}
	export class SignupResult {
	    userId: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new SignupResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userId = source["userId"];
	        this.error = source["error"];
	    }
	}

}

