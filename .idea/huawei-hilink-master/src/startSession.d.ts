export type SessionData = {
    TokInfo: string;
    SesInfo: string;
    url: string;
};
export declare function hilinkLogin(url: string): Promise<any>;
export declare function logout(url: string): Promise<void>;
export declare function login(url: string, user: string, password: string | undefined): Promise<void>;
export declare function getToken(url: string, cooke: string): Promise<any>;
export declare function startSession(url: string): Promise<SessionData>;
//# sourceMappingURL=startSession.d.ts.map