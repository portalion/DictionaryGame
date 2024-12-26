import { useReducer } from 'react';
import RoomSelector from './components/RoomSelector';
import Room from './components/Room';
import { AppState, GlobalData, GlobalStateAction } from './App.types';

function globalStateReducer(state: GlobalData, action: GlobalStateAction): GlobalData {
  let result = state;
  switch (action.type) {
    case 'ChangeGlobalData':
      result = { ...state, ...action.newData };
      break;
    case 'ChangeRoomCode':
      result = { ...state, currentRoomId: action.newRoomCode };
      break;
    case 'ChangeRoomState':
      result = { ...state, state: action.newState };
      break;
    case 'ChangeUsername':
      result = { ...state, username: action.newUsername };
      break;
  }
  if (result.currentRoomId != '') result.state = AppState.InRoom;
  else result.state = AppState.Idle;

  return result;
}

function App() {
  const [appState, dispatch] = useReducer(globalStateReducer, new GlobalData());

  switch (appState.state) {
    case AppState.Idle:
      return (
        <RoomSelector
          dispatch={dispatch}
          globalState={appState}
        />
      );
    case AppState.InRoom:
      return (
        <Room
          dispatch={dispatch}
          globalState={appState}
        />
      );
  }
}

export default App;
