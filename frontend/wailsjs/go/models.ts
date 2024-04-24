export namespace nav {
	
	export class FileInfo {
	    name: string;
	    rawSize: number;
	    isDir: boolean;
	    mode: string;
	    username: string;
	    groupName: string;
	    size: string;
	    modTime: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.rawSize = source["rawSize"];
	        this.isDir = source["isDir"];
	        this.mode = source["mode"];
	        this.username = source["username"];
	        this.groupName = source["groupName"];
	        this.size = source["size"];
	        this.modTime = source["modTime"];
	    }
	}
	export class PreviewInfo {
	    dirFiles?: FileInfo[];
	    utf8?: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new PreviewInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirFiles = this.convertValues(source["dirFiles"], FileInfo);
	        this.utf8 = source["utf8"];
	        this.type = source["type"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

