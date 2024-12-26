export enum AppState {
    Idle = 0,
    InRoom = 1,
}

export class GlobalData {
    state: AppState = AppState.Idle;
    username: string = '';
    currentRoomId: string = '';
}

export type GlobalStateAction =
    | { type: 'ChangeUsername'; newUsername: string }
    | { type: 'ChangeRoomCode'; newRoomCode: string }
    | { type: 'ChangeGlobalData'; newData: Partial<Omit<GlobalData, 'state'>> }
    | { type: 'ChangeRoomState'; newState: AppState };
