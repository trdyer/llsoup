import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';
import { AllData, Thermometer } from './interfaces';
@Injectable({
  providedIn: 'root',
})
export class ThermometerService {
  constructor(private http: HttpClient) {}

  public GetAllData(): Observable<Thermometer[]> {
    return this.http.get<AllData>('//thermometer.tristandyer.ca/api/all').pipe(
      map(d =>
        Object.keys(d).map(
          t =>
            ({
              ...d[t],
              city: t,
            }) as Thermometer,
        ),
      ),
    );
  }
}
