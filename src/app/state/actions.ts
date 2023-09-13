import { createAction, props } from '@ngrx/store';
import { Thermometer } from '../interfaces';

export const getThermometers = createAction('Thermometers] Load data');

export const getThermometersSuccess = createAction('Thermometers] Load data success', props<{ data: Thermometer[] }>());

export const getThermometersFailure = createAction('Thermometers] Load data failure', props<{ error: string }>());
