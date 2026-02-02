export namespace exif {
	
	export class ExifData {
	    // Go type: time
	    dateTime: any;
	    camera?: string;
	    lens?: string;
	    iso?: number;
	    fNumber?: string;
	    exposureTime?: string;
	    focalLength?: string;
	    width?: number;
	    height?: number;
	    dimensions?: string;
	
	    static createFrom(source: any = {}) {
	        return new ExifData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dateTime = this.convertValues(source["dateTime"], null);
	        this.camera = source["camera"];
	        this.lens = source["lens"];
	        this.iso = source["iso"];
	        this.fNumber = source["fNumber"];
	        this.exposureTime = source["exposureTime"];
	        this.focalLength = source["focalLength"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.dimensions = source["dimensions"];
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

export namespace importer {
	
	export class VideoMetadata {
	    duration: number;
	    width: number;
	    height: number;
	    frameRate: string;
	    resolution: string;
	
	    static createFrom(source: any = {}) {
	        return new VideoMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.duration = source["duration"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.frameRate = source["frameRate"];
	        this.resolution = source["resolution"];
	    }
	}
	export class FileInfo {
	    name: string;
	    path: string;
	    size: number;
	    date: string;
	    type: string;
	    thumbnail?: string;
	    exif?: exif.ExifData;
	    videoMeta?: VideoMetadata;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.date = source["date"];
	        this.type = source["type"];
	        this.thumbnail = source["thumbnail"];
	        this.exif = this.convertValues(source["exif"], exif.ExifData);
	        this.videoMeta = this.convertValues(source["videoMeta"], VideoMetadata);
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
	export class ImportProgress {
	    current: number;
	    total: number;
	    percentage: number;
	    currentFile: string;
	    status: string;
	    successCount: number;
	    failureCount: number;
	    destinationPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current = source["current"];
	        this.total = source["total"];
	        this.percentage = source["percentage"];
	        this.currentFile = source["currentFile"];
	        this.status = source["status"];
	        this.successCount = source["successCount"];
	        this.failureCount = source["failureCount"];
	        this.destinationPath = source["destinationPath"];
	    }
	}

}

export namespace main {
	
	export class Config {
	    source: string;
	    destination: string;
	    imgFormats: string[];
	    videoFormats: string[];
	    rawFormats: string[];
	    imgRelativePath: string;
	    videoRelativePath: string;
	    rawRelativePath: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.source = source["source"];
	        this.destination = source["destination"];
	        this.imgFormats = source["imgFormats"];
	        this.videoFormats = source["videoFormats"];
	        this.rawFormats = source["rawFormats"];
	        this.imgRelativePath = source["imgRelativePath"];
	        this.videoRelativePath = source["videoRelativePath"];
	        this.rawRelativePath = source["rawRelativePath"];
	    }
	}

}

