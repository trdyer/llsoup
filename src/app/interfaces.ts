export interface Thermometer {
  raised: number;
  goal: number;
  city: string;
}

export interface AllData {
  [key: string]: Thermometer;
}
