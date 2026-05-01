import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import {
  ActorService,
  BfsResponse,
  Bfs8Response,
  PathItem,
} from './actor.service';

type Mode = 'bfs' | 'bfs8' | null;

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App implements OnInit {
  private actorService = inject(ActorService);

  readonly actors = signal<string[]>([]);
  readonly actorsCount = signal<number>(0);
  readonly loadingActors = signal<boolean>(true);
  readonly loadError = signal<string | null>(null);

  readonly from = signal<string>('');
  readonly to = signal<string>('');

  readonly searching = signal<boolean>(false);
  readonly mode = signal<Mode>(null);

  readonly bfsResult = signal<BfsResponse | null>(null);
  readonly bfs8Result = signal<Bfs8Response | null>(null);
  readonly errorMessage = signal<string | null>(null);

  readonly canSearch = computed(
    () =>
      !this.searching() &&
      this.from().trim().length > 0 &&
      this.to().trim().length > 0,
  );

  ngOnInit(): void {
    this.actorService.getActors().subscribe({
      next: (res) => {
        this.actors.set(res.actors);
        this.actorsCount.set(res.count);
        this.loadingActors.set(false);
      },
      error: (err) => {
        this.loadingActors.set(false);
        this.loadError.set(
          'Não foi possível carregar a lista de atores. Verifique se o backend está rodando em http://localhost:8081',
        );
        console.error(err);
      },
    });
  }

  swap(): void {
    const a = this.from();
    this.from.set(this.to());
    this.to.set(a);
  }

  clearResults(): void {
    this.bfsResult.set(null);
    this.bfs8Result.set(null);
    this.errorMessage.set(null);
    this.mode.set(null);
  }

  runBfs(): void {
    if (!this.canSearch()) return;
    this.clearResults();
    this.searching.set(true);
    this.mode.set('bfs');
    this.actorService.bfs(this.from().trim(), this.to().trim()).subscribe({
      next: (res) => {
        this.bfsResult.set(res);
        this.searching.set(false);
      },
      error: (err) => {
        this.searching.set(false);
        if (err.status === 404) {
          this.errorMessage.set('Ator não encontrado no grafo.');
        } else {
          this.errorMessage.set('Erro ao consultar o backend.');
        }
        console.error(err);
      },
    });
  }

  runBfs8(): void {
    if (!this.canSearch()) return;
    this.clearResults();
    this.searching.set(true);
    this.mode.set('bfs8');
    this.actorService.bfs8(this.from().trim(), this.to().trim(), 8).subscribe({
      next: (res) => {
        this.bfs8Result.set(res);
        this.searching.set(false);
      },
      error: (err) => {
        this.searching.set(false);
        if (err.status === 404) {
          this.errorMessage.set('Ator não encontrado no grafo.');
        } else {
          this.errorMessage.set('Erro ao consultar o backend.');
        }
        console.error(err);
      },
    });
  }

  isMovieNode(_node: string, index: number): boolean {
    return index % 2 === 1;
  }

  trackPath(_index: number, item: PathItem): string {
    return item.path.join('|');
  }
}
