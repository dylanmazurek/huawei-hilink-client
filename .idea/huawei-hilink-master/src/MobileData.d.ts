import { SessionData } from "./startSession";
import { ExportFormat } from "./utils/Constants";
type MobileStatus = 'on' | 'off';
export declare function controlMobileData(sessionData: SessionData, mobileStatus: MobileStatus): Promise<void>;
export declare function reconnect(sessionData: SessionData): Promise<void>;
export declare function status(sessionData: SessionData, exportFile: string, exportFormat: ExportFormat): Promise<void>;
export {};
//# sourceMappingURL=MobileData.d.ts.map