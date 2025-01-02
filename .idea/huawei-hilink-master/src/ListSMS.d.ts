import { SessionData } from './startSession';
import { ExportFormat } from "./utils/Constants";
export declare function getSMSByUsers(sessionData: SessionData, phone: string, pageindex: number, exportFile: string, exportFormat: ExportFormat, deleteAfter: boolean): Promise<any>;
export declare function getContactSMSPages(sessionData: SessionData, phone: string, exportFile: string, exportFormat: ExportFormat): Promise<number>;
export declare function getSMSPages(sessionData: SessionData, exportFile: string, exportFormat: ExportFormat): Promise<number>;
export declare function deleteMessage(sessionData: SessionData, messageId: string): Promise<void>;
export declare function sendMessage(sessionData: SessionData, phones: string, message: string): Promise<void>;
export declare function getInBoxSMS(sessionData: SessionData, deleteAfter: boolean, exportFile: string, exportFormat: ExportFormat): Promise<any>;
export declare function getSMSContacts(sessionData0: SessionData, pageindex: number, exportFile: string, exportFormat: ExportFormat): Promise<any>;
//# sourceMappingURL=ListSMS.d.ts.map