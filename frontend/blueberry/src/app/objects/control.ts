export interface Room {
  id: string;
  name: string;
  controlGroups: Object; // map <string, ControlGroup>
  selectedControlGroup: string;
}

export function isRoom(o: Object): o is Room {
  return (
    o &&
    o.hasOwnProperty("controlGroups") &&
    o.hasOwnProperty("id") &&
    o.hasOwnProperty("name")
  );
}

export interface ControlGroup {
  id: string;
  name: string;
  displayBlocks: DisplayBlock[];
  inputs: Input[];
  audioGroups: AudioGroup[];
  presentGroups: PresentGroup[];
  support: Support;
  level: number;
  muted: boolean;
  screens: string[];
  powerOn: boolean;

  // public getAudioDevice(cg: ControlGroup, id: string): AudioDevice {
  //     for (const g of cg.audioGroups) {
  //         for (const device of g.audioDevices) {
  //             if (device.id === id) {
  //                 return device;
  //             }
  //         }
  //     }
  // }
}

export interface Support {
  helpRequested: boolean;
  helpMessage: string;
  helpEnabled: boolean;
}

export interface DisplayBlock {
  id: string;
  outputs: IconPair[];
  input: string;
  blanked: boolean;
  shareOptions: string[];

  // getOutputNameList(): string[] {
  //     const toReturn: string[] = [];
  //     for (const o of this.outputs) {
  //         toReturn.push(o.name);
  //     }
  //     return toReturn;
  // }
}

export interface Input {
  id: string;
  name: string;
  icon: string;
  subInputs: Input[];
  disabled: boolean;
}

export interface AudioGroup {
  id: string;
  name: string;
  audioDevices: AudioDevice[];
  muted: boolean;
}

export interface AudioDevice {
  id: string;
  name: string;
  icon: string;
  level: number;
  muted: boolean;
}

export interface PresentGroup {
  id: string;
  name: string;
  items: PresentItem[];
}

export interface PresentItem {
  id: string;
  name: string;
}

export interface IconPair {
  id: string;
  icon: string;
  name: string;
}

export const CONTROL_TAB = "Control";
export const AUDIO_TAB = "Audio";
export const PRESENT_TAB = "Present";
export const HELP_TAB = "Help";
