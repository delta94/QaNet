import React from 'react';
import Link from 'next/link';
import styled from 'styled-components';

const A = styled('a')`
  color: #fff;
  &:hover {
    color: #fff;
  }
`;

const CustomA = React.forwardRef(({ onClick, href, text }: any, ref: any) => (
  <A className="ui blue button" href={href} onClick={onClick} ref={ref}>
    {text}
  </A>
));

interface IProps {
  href: string;
  text: string;
}

const MenuButton = function({ href, text }: IProps) {
  return (
    <Link href={href}>
      <CustomA text={text} />
    </Link>
  );
};

export default MenuButton;
