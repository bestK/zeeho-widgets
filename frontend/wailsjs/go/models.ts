export namespace main {
	
	export class Config {
	    token: string;
	    vehicleId: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.vehicleId = source["vehicleId"];
	    }
	}
	export class LocationData {
	    coordinateSystem: string;
	    latitude: number;
	    longitude: number;
	    locationTime: string;
	    altitude: number;
	    address?: string;
	
	    static createFrom(source: any = {}) {
	        return new LocationData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.coordinateSystem = source["coordinateSystem"];
	        this.latitude = source["latitude"];
	        this.longitude = source["longitude"];
	        this.locationTime = source["locationTime"];
	        this.altitude = source["altitude"];
	        this.address = source["address"];
	    }
	}
	export class VehicleData {
	    vehicleName: string;
	    bmssoc: string;
	    hmiRidableMile: string;
	    vehiclePicUrl: string;
	    location?: LocationData;
	
	    static createFrom(source: any = {}) {
	        return new VehicleData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicleName = source["vehicleName"];
	        this.bmssoc = source["bmssoc"];
	        this.hmiRidableMile = source["hmiRidableMile"];
	        this.vehiclePicUrl = source["vehiclePicUrl"];
	        this.location = this.convertValues(source["location"], LocationData);
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

