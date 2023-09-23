import { Component, OnInit } from '@angular/core';
import { Store } from '@ngrx/store';
import { Observable, Subject, merge, of, scan } from 'rxjs';
import { Thermometer } from './interfaces';
import { getThermometers } from './state/actions';
import { selectThermometers, selectThermometersLoading } from './state/reducer';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  public title = 'Thermometer Data';

  public allData$!: Observable<Thermometer[]>;
  public loading$!: Observable<boolean>;
  public listView$: Observable<boolean>;
  private listViewSubject$: Subject<void>;

  constructor(private store: Store) {
    this.listViewSubject$ = new Subject<void>();
    this.listView$ = merge(of(true), this.listViewSubject$).pipe(scan((acc, _) => !acc, true));
  }

  public ngOnInit() {
    this.store.dispatch(getThermometers());

    this.allData$ = this.store.select(selectThermometers);
    this.loading$ = this.store.select(selectThermometersLoading);
  }

  public toggleListView(): void {
    this.listViewSubject$.next();
  }

  public refresh(): void {
    this.store.dispatch(getThermometers());
  }
}
