export interface Thermometer {
  raised: number;
  goal: number;
}

export interface AllData {
  halifax: Thermometer;
  calgary: Thermometer;
  edmonton: Thermometer;
  montreal: Thermometer;
  stjohns: Thermometer;
  toronto: Thermometer;
  vancouver: Thermometer;
  london: Thermometer;
  ottawa: Thermometer;
  winnipeg: Thermometer;
}
