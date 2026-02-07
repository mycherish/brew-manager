export namespace main {
	
	export class ActionResponse {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ActionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class BrewPackage {
	    name: string;
	    version: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new BrewPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.status = source["status"];
	    }
	}
	export class BrewData {
	    formulae: BrewPackage[];
	    casks: BrewPackage[];
	
	    static createFrom(source: any = {}) {
	        return new BrewData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.formulae = this.convertValues(source["formulae"], BrewPackage);
	        this.casks = this.convertValues(source["casks"], BrewPackage);
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

}

