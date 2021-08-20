import { TestBed } from '@angular/core/testing';

import { ThermometerService } from './thermometer.service';

describe('ThermometerService', () => {
  let service: ThermometerService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ThermometerService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
