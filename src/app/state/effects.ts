import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap } from 'rxjs';
import { ThermometerService } from '../thermometer.service';
import { getThermometers, getThermometersFailure, getThermometersSuccess } from './actions';

@Injectable()
export class ThermometerEffects {
  public loadThermometers$ = createEffect(() =>
    this.actions$.pipe(
      ofType(getThermometers),
      switchMap(() =>
        this.thermSvc.GetAllData().pipe(
          map(data => getThermometersSuccess({ data })),
          catchError((err: HttpErrorResponse) => of(getThermometersFailure({ error: err.message }))),
        ),
      ),
    ),
  );

  constructor(
    private actions$: Actions,
    private thermSvc: ThermometerService,
  ) {}
}
