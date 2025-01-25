Some ideas:

Command format:

```
{ action: 'registerShip', name: string }

{ action: 'forward' }

{ action: 'stop' }

{ action: 'shoot' }

{ action: 'rotate', angle: number }
```

Error format:

```
{ error: string }
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
