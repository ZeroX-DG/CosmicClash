Some ideas:

Command format:

```
{ action: 'forward' }

{ action: 'stop' }

{ action: 'shoot' }

{ action: 'rotate', angle: number }
```

Game state format:

```
{
    ships: [
        {
            name: string;
            location: [number, number],
            health: number,
        }
    ]
}
```
