import { Type } from "serializer.ts/Decorators";
import { EventEmitter } from "@angular/core";
import { AudioConfiguration } from "./objects";
// import { APIService } from "../services/api.service";

export const POWER = "power";
export const INPUT = "input";
export const BLANKED = "blanked";
export const MUTED = "muted";
export const VOLUME = "volume";

export class Device {
  name: string;
  displayname: string;
  icon: string;

  constructor(name: string, displayname: string, icon: string) {
    this.name = name;
    this.displayname = displayname;
    this.icon = icon;
  }

  public static filterDevices<T extends Device>(
    names: string[],
    devices: T[]
  ): T[] {
    if (names == null || devices == null) {
      return [];
    }
    const ret: T[] = [];

    for (const name of names) {
      const dev = devices.find(d => d.name === name);
      if (dev != null) {
        ret.push(dev);
      }
    }

    return ret;
  }

  public static getDeviceByName<T extends Device>(
    name: string,
    devices: T[]
  ): T {
    if (name == null || devices == null) {
      return null;
    }

    for (const d of devices) {
      if (d.name === name) {
        return d;
      }
    }

    return null;
  }

  public getName(): string {
    return this.name;
  }

  public getDisplayName(): string {
    return this.displayname;
  }

  public getIcon(): string {
    return this.icon;
  }
}

export class Input extends Device {
  click: EventEmitter<null> = new EventEmitter();
  subInputs: Input[] = [];

  constructor(name: string, displayname: string, icon: string, subs: Input[]) {
    super(name, displayname, icon);
    this.subInputs = subs;
  }

  public static getInput(name: string, inputs: Input[]): Input {
    for (const i of inputs) {
      if (i.name === name) {
        return i;
      }
      if (i.subInputs !== undefined && i.subInputs.length > 0) {
        for (const sub of i.subInputs) {
          if (sub.name === name) {
            return sub;
          }
        }
      }
    }
  }
}

export class Output extends Device {
  power: string;
  input: Input;

  powerEmitter: EventEmitter<string>;

  public static getPower(outputs: Output[]): string {
    for (const o of outputs) {
      if (o.power === "on") {
        return o.power;
      }
    }

    return "standby";
  }

  public static isPoweredOn(outputs: Output[]): boolean {
    for (const o of outputs) {
      if (o.power !== "on") {
        return false;
      }
    }

    return true;
  }

  public static getInput(outputs: Output[]): Input {
    let input: Input = null;

    for (const o of outputs) {
      if (input == null) {
        input = o.input;
      } else if (o.input !== input) {
        // this means the input that appears selected may not actually be selected on all displays.
        // to get the ~correct~ behavior, return null.
        return o.input;
      }
    }

    return input;
  }

  public static setPower(s: string, outputs: Output[]) {
    outputs.forEach(o => (o.power = s));
  }

  public static setInput(i: Input, outputs: Output[]) {
    outputs.forEach(o => (o.input = i));
  }

  constructor(
    name: string,
    displayname: string,
    power: string,
    input: Input,
    icon: string
  ) {
    super(name, displayname, icon);
    this.power = power;
    this.input = input;

    this.powerEmitter = new EventEmitter();
  }

  public getInputIcon(): string {
    if (this.input == null) {
      return this.icon;
    }

    return this.input.icon;
  }
}

export class Display extends Output {
  blanked: boolean;

  public static getDisplayListFromNames(
    names: string[],
    displaysSource: Display[]
  ): Display[] {
    return displaysSource.filter(d => names.includes(d.name));
  }

  constructor(
    name: string,
    displayname: string,
    power: string,
    input: Input,
    blanked: boolean,
    icon: string
  ) {
    super(name, displayname, power, input, icon);
    this.blanked = blanked;
  }

  // returns true iff all are blanked
  public static getBlank(displays: Display[]): boolean {
    for (const d of displays) {
      if (!d.blanked) {
        return false;
      }
    }

    return true;
  }

  public static setBlank(b: boolean, displays: Display[]) {
    displays.forEach(d => (d.blanked = b));
  }

  // public getAudioConfiguration(): AudioConfiguration {
  //   return APIService.room.uiconfig.audioConfiguration.find(
  //     a => a.display === this.name
  //   );
  // }
}

export class AudioDevice extends Output {
  muted: boolean;
  volume: number;
  type: string;

  mixlevel: number;
  mixmute: boolean;

  constructor(
    name: string,
    displayname: string,
    power: string,
    input: Input,
    muted: boolean,
    volume: number,
    icon: string,
    type: string,
    mixlevel: number
  ) {
    super(name, displayname, power, input, icon);
    this.muted = muted;
    this.volume = volume;
    this.type = type;
    this.mixlevel = mixlevel;
  }

  // return average of all volumes
  public static getVolume(audioDevices: AudioDevice[]): number {
    if (audioDevices == null) {
      return 0;
    }

    let volume = 0;

    audioDevices.forEach(a => (volume += a.volume));

    return volume / audioDevices.length;
  }

  // returns true iff both are muted
  public static getMute(audioDevices: AudioDevice[]): boolean {
    if (audioDevices == null) {
      return false;
    }

    for (const a of audioDevices) {
      if (!a.muted) {
        return false;
      }
    }

    return true;
  }

  public static setVolume(v: number, audioDevices: AudioDevice[]) {
    audioDevices.forEach(a => (a.volume = v));
  }

  public static setMute(m: boolean, audioDevices: AudioDevice[]) {
    audioDevices.forEach(a => (a.muted = m));
  }
}
