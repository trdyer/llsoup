import { Component, Input } from '@angular/core';
import { Thermometer } from '../interfaces';

@Component({
  selector: 'app-thermometer',
  templateUrl: './thermometer.component.html',
  styleUrls: ['./thermometer.component.scss'],
})
export class ThermometerComponent {
  @Input() data!: Thermometer;
}
