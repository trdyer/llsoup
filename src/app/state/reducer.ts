import { Action, createFeatureSelector, createReducer, createSelector, on } from '@ngrx/store';
import { Thermometer } from '../interfaces';
import { getThermometers, getThermometersFailure, getThermometersSuccess } from './actions';

export interface ThermometerState {
  thermometers: Thermometer[];
  loading: boolean;
  error: string;
}
export const initialState: ThermometerState = {
  thermometers: [],
  loading: false,
  error: '',
};

export const thermometerReducer = createReducer(
  initialState,
  on(getThermometers, state => ({
    ...state,
    loading: true,
    error: '',
  })),
  on(getThermometersSuccess, (state, action) => ({
    ...state,
    loading: false,
    error: '',
    thermometers: action.data,
  })),
  on(getThermometersFailure, (state, action) => ({
    ...state,
    loading: false,
    error: action.error,
    thermometers: [],
  }))
);

export function reducer(state: ThermometerState | undefined, action: Action): ThermometerState {
  return thermometerReducer(state, action);
}
export const selectRoot = createFeatureSelector<ThermometerState>('root');

export const selectThermometers = createSelector(selectRoot, (state: ThermometerState) => state.thermometers);
export const selectThermometersLoading = createSelector(selectRoot, (state: ThermometerState) => state.loading);
export const selectThermometersError = createSelector(selectRoot, (state: ThermometerState) => state.error);
