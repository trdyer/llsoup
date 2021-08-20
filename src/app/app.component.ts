import { ChangeDetectionStrategy, Component } from "@angular/core";
import { Observable } from "rxjs";
import { AllData, Thermometer } from "./interfaces";
import { ThermometerService } from "./thermometer.service";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AppComponent {
  title = "Thermometer Data";

  public allData$: Observable<AllData>;
  constructor(private thermSvc: ThermometerService) {
    this.allData$ = this.thermSvc.GetAllData();
  }

  public getKeys(data: AllData): string[] {
    return Object.keys(data);
  }
}
