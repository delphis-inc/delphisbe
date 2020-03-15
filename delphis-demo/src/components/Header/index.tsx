import React from 'react';
import './style.scss';

export interface Props {}

function Header(props: Props) {
    return (
        <div className="delphis-header w-full flex justify-center">
            <div className="pl-2 pr-2 px-4 py-2 mx-6 w-full md:w-8/12 lg:w-8/12 xl:w-8/12 mx-6">
                <div className="flex justify-between items-center">
                    <a className="w-1/6 sm:w-1/6 md:w-1/6 lg:w-1/6 xl:w-1/6 logo" href="/"><img src="/logo.png" alt=""/></a>
                    <div className="links">
                        <a className="mr-6" href="/about">Dafuq?</a>
                        <a href="/contact">Contact</a>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Header