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
}
