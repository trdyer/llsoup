import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { Observable } from "rxjs";
import { AllData } from "./interfaces";

@Injectable({
  providedIn: "root",
})
export class ThermometerService {
  constructor(private http: HttpClient) {}

  public GetAllData(): Observable<AllData> {
    return this.http.get<AllData>(
      "https://safe-reef-59505.herokuapp.com/api/all"
    );
  }
}
