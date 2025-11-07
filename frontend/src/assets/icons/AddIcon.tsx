import React from 'react';
import type { SVGProps } from 'react';

export function AddIcon(props: SVGProps<SVGSVGElement>) {
    return (<svg xmlns="http://www.w3.org/2000/svg" width={24} height={24} viewBox="0 0 24 24" {...props}><path fill="#7d82b8" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6z"></path></svg>);
}

export default AddIcon;