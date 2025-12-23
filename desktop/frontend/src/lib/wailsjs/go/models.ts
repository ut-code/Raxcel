export namespace main {
  export class ChatResult {
    ok: boolean;
    message: string;

    static createFrom(source: any = {}) {
      return new ChatResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.ok = source["ok"];
      this.message = source["message"];
    }
  }
  export class CheckResult {
    ok: boolean;
    userId?: string;
    error?: string;

    static createFrom(source: any = {}) {
      return new CheckResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.ok = source["ok"];
      this.userId = source["userId"];
      this.error = source["error"];
    }
  }
  export class MessageItem {
    id: string;
    userId: string;
    content: string;
    role: string;
    createdAt: string;

    static createFrom(source: any = {}) {
      return new MessageItem(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.id = source["id"];
      this.userId = source["userId"];
      this.content = source["content"];
      this.role = source["role"];
      this.createdAt = source["createdAt"];
    }
  }
  export class GetMessagesResult {
    ok: boolean;
    messages: MessageItem[];
    error?: string;

    static createFrom(source: any = {}) {
      return new GetMessagesResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.ok = source["ok"];
      this.messages = this.convertValues(source["messages"], MessageItem);
      this.error = source["error"];
    }

    convertValues(a: any, classs: any, asMap: boolean = false): any {
      if (!a) {
        return a;
      }
      if (a.slice && a.map) {
        return (a as any[]).map((elem) => this.convertValues(elem, classs));
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
  export class LoginResult {
    ok: boolean;
    message: string;
    token?: string;

    static createFrom(source: any = {}) {
      return new LoginResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.ok = source["ok"];
      this.message = source["message"];
      this.token = source["token"];
    }
  }

  export class RegisterResult {
    ok: boolean;
    message: string;
    userId?: string;

    static createFrom(source: any = {}) {
      return new RegisterResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.ok = source["ok"];
      this.message = source["message"];
      this.userId = source["userId"];
    }
  }
  export class SignOutResult {
    ok: boolean;
    error?: string;

    static createFrom(source: any = {}) {
      return new SignOutResult(source);
    }

    constructor(source: any = {}) {
      if ("string" === typeof source) source = JSON.parse(source);
      this.ok = source["ok"];
      this.error = source["error"];
    }
  }
}
