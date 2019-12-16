export class Room {
  id: string;
  name: string;
  controlGroups: Map<string, ControlGroup>;
  selectedControlGroup: string;
}

export class ControlGroup {
  id: string;
  name: string;
  displays: Display[];
  inputs: Input[];
  audioGroups: AudioGroup[];
  presentGroups: PresentGroup[];
  support: Support;
  level: number;
  muted: boolean;

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

export class Support {
  helpRequested: boolean;
  helpMessage: string;
  helpEnabled: boolean;
}

export class Display {
  id: string;
  outputs: IconPair[];
  input: string;
  blanked: boolean;

  // getOutputNameList(): string[] {
  //     const toReturn: string[] = [];
  //     for (const o of this.outputs) {
  //         toReturn.push(o.name);
  //     }
  //     return toReturn;
  // }
}

export class Input {
  id: string;
  name: string;
  icon: string;
  subInputs: Input[];
  disabled: boolean;
}

export class AudioGroup {
  id: string;
  name: string;
  audioDevices: AudioDevice[];
  muted: boolean;
}

export class AudioDevice {
  id: string;
  name: string;
  icon: string;
  level: number;
  muted: boolean;
}

export class PresentGroup {
  id: string;
  name: string;
  items: PresentItem[];
}

export class PresentItem {
  id: string;
  name: string;
}

export class IconPair {
  id: string;
  icon: string;
  name: string;
}

export const CONTROL_TAB = "Control";
export const AUDIO_TAB = "Audio";
export const PRESENT_TAB = "Present";
export const HELP_TAB = "Help";
