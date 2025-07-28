export namespace main {
	
	export class Config {
	    token: string;
	    vehicleId: string;
	    updateInterval: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.vehicleId = source["vehicleId"];
	        this.updateInterval = source["updateInterval"];
	    }
	}
	export class EncryptInfo {
	    key: string;
	    iv: string;
	    encryptValue: string;
	
	    static createFrom(source: any = {}) {
	        return new EncryptInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.iv = source["iv"];
	        this.encryptValue = source["encryptValue"];
	    }
	}
	export class IotProperty {
	    name: string;
	    identify: string;
	    value: string;
	    time: string;
	    dbUpdateTime?: string;
	    describe?: string;
	
	    static createFrom(source: any = {}) {
	        return new IotProperty(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.identify = source["identify"];
	        this.value = source["value"];
	        this.time = source["time"];
	        this.dbUpdateTime = source["dbUpdateTime"];
	        this.describe = source["describe"];
	    }
	}
	export class Location {
	    longitude: number;
	    latitude: number;
	    altitude: number;
	    coordinateSystem: string;
	    locationTime: string;
	    address?: string;
	
	    static createFrom(source: any = {}) {
	        return new Location(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.longitude = source["longitude"];
	        this.latitude = source["latitude"];
	        this.altitude = source["altitude"];
	        this.coordinateSystem = source["coordinateSystem"];
	        this.locationTime = source["locationTime"];
	        this.address = source["address"];
	    }
	}
	export class VehicleData {
	    vinNo: string;
	    deviceName: string;
	    bindStartTime: string;
	    fate: string;
	    vehicleName: string;
	    vehiclePicUrl: string;
	    vehicleBackPicUrl: string;
	    shareName?: string;
	    isShared: string;
	    rideMileageMonth: string;
	    ridingTimeMonth: string;
	    ridingTimeMonthUnitMinute: string;
	    avgVelocityMonth: string;
	    bmssoc: string;
	    hmiRidableMile: string;
	    gsmRxLev: string;
	    gsmRxLevValue: string;
	    bluetoothAddress: string;
	    hmiBluetoothAddress: string;
	    chargeState: string;
	    fullChargeTime: string;
	    pressure: string;
	    pressureValue: string;
	    headLockState: string;
	    rideState: string;
	    greenContribution: string;
	    otaVersion: string;
	    shareUserId?: string;
	    bindingUserId: number;
	    carMaster: string;
	    vehicleType: string;
	    vehicleTypeName: string;
	    encryptInfo: EncryptInfo;
	    redPoint: number;
	    hmiRidableMileAbnormalShow?: string;
	    expectFullTimeDescribe?: string;
	    location: Location;
	    navigationType: number;
	    navigation: string;
	    projection: string;
	    motoPlay: number;
	    wifiAddress: string;
	    bluetoothSearch: boolean;
	    vehicleTypeDetailId: number;
	    shareEndTime?: string;
	    residualSeconds?: string;
	    supportNetworkUnlock: number;
	    intelligentType?: string;
	    totalRideMile: string;
	    maxMileage: string;
	    deviceType: number;
	    broadcastType: string;
	    supportUnlock: number;
	    bindDate?: string;
	    cyclingEventStatisticFlag: boolean;
	    whetherChargeState: boolean;
	    firstBindDate: string;
	    maxRangeMileage: string;
	    onlineStatus: string;
	    activationDate: string;
	    lastUseDate: number;
	    rechargeEndDate: string;
	    openCushionFlag: boolean;
	    openStorageBoxFlag: boolean;
	    loudlySearchCar: number;
	    mmiUuid: string;
	    productKey: string;
	    iotInstanceId: string;
	    iotProperties: IotProperty[];
	    serviceRechargeStatus: string;
	    refreshTime: string;
	    vehicleScalePicUrl: string;
	    gaodeLincenseVinNo: string;
	    gaodeLincenseId: string;
	
	    static createFrom(source: any = {}) {
	        return new VehicleData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vinNo = source["vinNo"];
	        this.deviceName = source["deviceName"];
	        this.bindStartTime = source["bindStartTime"];
	        this.fate = source["fate"];
	        this.vehicleName = source["vehicleName"];
	        this.vehiclePicUrl = source["vehiclePicUrl"];
	        this.vehicleBackPicUrl = source["vehicleBackPicUrl"];
	        this.shareName = source["shareName"];
	        this.isShared = source["isShared"];
	        this.rideMileageMonth = source["rideMileageMonth"];
	        this.ridingTimeMonth = source["ridingTimeMonth"];
	        this.ridingTimeMonthUnitMinute = source["ridingTimeMonthUnitMinute"];
	        this.avgVelocityMonth = source["avgVelocityMonth"];
	        this.bmssoc = source["bmssoc"];
	        this.hmiRidableMile = source["hmiRidableMile"];
	        this.gsmRxLev = source["gsmRxLev"];
	        this.gsmRxLevValue = source["gsmRxLevValue"];
	        this.bluetoothAddress = source["bluetoothAddress"];
	        this.hmiBluetoothAddress = source["hmiBluetoothAddress"];
	        this.chargeState = source["chargeState"];
	        this.fullChargeTime = source["fullChargeTime"];
	        this.pressure = source["pressure"];
	        this.pressureValue = source["pressureValue"];
	        this.headLockState = source["headLockState"];
	        this.rideState = source["rideState"];
	        this.greenContribution = source["greenContribution"];
	        this.otaVersion = source["otaVersion"];
	        this.shareUserId = source["shareUserId"];
	        this.bindingUserId = source["bindingUserId"];
	        this.carMaster = source["carMaster"];
	        this.vehicleType = source["vehicleType"];
	        this.vehicleTypeName = source["vehicleTypeName"];
	        this.encryptInfo = this.convertValues(source["encryptInfo"], EncryptInfo);
	        this.redPoint = source["redPoint"];
	        this.hmiRidableMileAbnormalShow = source["hmiRidableMileAbnormalShow"];
	        this.expectFullTimeDescribe = source["expectFullTimeDescribe"];
	        this.location = this.convertValues(source["location"], Location);
	        this.navigationType = source["navigationType"];
	        this.navigation = source["navigation"];
	        this.projection = source["projection"];
	        this.motoPlay = source["motoPlay"];
	        this.wifiAddress = source["wifiAddress"];
	        this.bluetoothSearch = source["bluetoothSearch"];
	        this.vehicleTypeDetailId = source["vehicleTypeDetailId"];
	        this.shareEndTime = source["shareEndTime"];
	        this.residualSeconds = source["residualSeconds"];
	        this.supportNetworkUnlock = source["supportNetworkUnlock"];
	        this.intelligentType = source["intelligentType"];
	        this.totalRideMile = source["totalRideMile"];
	        this.maxMileage = source["maxMileage"];
	        this.deviceType = source["deviceType"];
	        this.broadcastType = source["broadcastType"];
	        this.supportUnlock = source["supportUnlock"];
	        this.bindDate = source["bindDate"];
	        this.cyclingEventStatisticFlag = source["cyclingEventStatisticFlag"];
	        this.whetherChargeState = source["whetherChargeState"];
	        this.firstBindDate = source["firstBindDate"];
	        this.maxRangeMileage = source["maxRangeMileage"];
	        this.onlineStatus = source["onlineStatus"];
	        this.activationDate = source["activationDate"];
	        this.lastUseDate = source["lastUseDate"];
	        this.rechargeEndDate = source["rechargeEndDate"];
	        this.openCushionFlag = source["openCushionFlag"];
	        this.openStorageBoxFlag = source["openStorageBoxFlag"];
	        this.loudlySearchCar = source["loudlySearchCar"];
	        this.mmiUuid = source["mmiUuid"];
	        this.productKey = source["productKey"];
	        this.iotInstanceId = source["iotInstanceId"];
	        this.iotProperties = this.convertValues(source["iotProperties"], IotProperty);
	        this.serviceRechargeStatus = source["serviceRechargeStatus"];
	        this.refreshTime = source["refreshTime"];
	        this.vehicleScalePicUrl = source["vehicleScalePicUrl"];
	        this.gaodeLincenseVinNo = source["gaodeLincenseVinNo"];
	        this.gaodeLincenseId = source["gaodeLincenseId"];
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

